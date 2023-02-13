package storage

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var (
	ErrMissingEnv         = errors.New("missing GOOGLE_CLOUD_SERVICE_ACCOUNT env var")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type GCSRepository struct {
	Bucket string
	ctx    context.Context
	client *storage.Client
	bucket *storage.BucketHandle
}

func GetCredentials(ctx context.Context) (*google.Credentials, error) {
	keyBase64 := os.Getenv("GOOGLE_CLOUD_SERVICE_ACCOUNT")
	if len(keyBase64) == 0 {
		return nil, ErrMissingEnv
	}
	rawDecodedText, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	creds, err := google.CredentialsFromJSON(ctx, rawDecodedText, "https://www.googleapis.com/auth/devstorage.read_write")
	if err != nil {
		return nil, err
	}
	return creds, err
}
func NewGcsRepository(bucket string) *GCSRepository {
	ctx := context.Background()
	creds, err := GetCredentials(ctx)
	if err != nil {
		log.Fatal(err)
	}
	client, err := storage.NewClient(ctx, option.WithCredentials(creds))
	b := client.Bucket(bucket)
	return &GCSRepository{ctx: context.Background(), client: client, Bucket: bucket, bucket: b}
}
func (gcs *GCSRepository) Read(path string) ([]byte, error) {
	obj, err := gcs.bucket.Object(path).NewReader(gcs.ctx)
	if err != nil {
		return nil, err
	}

	defer obj.Close()
	slurp, err := ioutil.ReadAll(obj)
	if err != nil {
		log.Printf("readFile: unable to read data from file %q: %v\n", path, err)
		return nil, err
	}
	return slurp, nil
	//log.Printf(d.w, "%s\n", bytes.SplitN(slurp, []byte("\n"), 2)[0])
}

func (gcs *GCSRepository) ReadImage(path string) (image.Image, error) {
	dat, err := gcs.Read(path)
	if err != nil {
		log.Println("cannot read file at:", path)
		return nil, err
	}
	imgFile := bytes.NewReader(dat)
	var img image.Image
	splittedPath := strings.Split(path, ".")
	ext := splittedPath[len(splittedPath)-1]

	if ext == "jpg" || ext == "jpeg" {
		img, err = jpeg.Decode(imgFile)
	} else {
		img, err = png.Decode(imgFile)
	}

	if err != nil {
		return nil, err
	}

	return img, nil
}
func (gcs *GCSRepository) Download(path string, toFile string) ([]byte, error) {
	dat, err := gcs.Read(path)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(toFile, dat, 0644)
	if err != nil {
		return nil, err
	}
	return dat, nil
}

func (gcs *GCSRepository) WriteImage(data *image.RGBA, output string) error {
	ctx, cancel := context.WithTimeout(gcs.ctx, time.Second*50)
	defer cancel()
	log.Println("write to file:", output)
	o := gcs.client.Bucket(gcs.Bucket).Object(output)
	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to upload is aborted if the
	// object's generation number does not match your precondition.
	// For an object that does not yet exist, set the DoesNotExist precondition.
	o = o.If(storage.Conditions{DoesNotExist: true})
	// If the live object already exists in your bucket, set instead a
	// generation-match precondition using the live object's generation number.
	// attrs, err := o.Attrs(ctx)
	// if err != nil {
	// 	return fmt.Errorf("object.Attrs: %v", err)
	// }
	// o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})
	readers := new(bytes.Buffer)
	// Upload an object with storage.Writer.
	wc := o.NewWriter(ctx)

	err := png.Encode(readers, data)
	if err != nil {
		return err
	}
	if _, err := io.Copy(wc, readers); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil
}

func (gcs *GCSRepository) Write(data []byte, output string) error {
	ctx, cancel := context.WithTimeout(gcs.ctx, time.Second*50)
	defer cancel()
	log.Println("write to file:", output)
	o := gcs.client.Bucket(gcs.Bucket).Object(output)
	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to upload is aborted if the
	// object's generation number does not match your precondition.
	// For an object that does not yet exist, set the DoesNotExist precondition.
	o = o.If(storage.Conditions{DoesNotExist: true})
	// If the live object already exists in your bucket, set instead a
	// generation-match precondition using the live object's generation number.
	// attrs, err := o.Attrs(ctx)
	// if err != nil {
	// 	return fmt.Errorf("object.Attrs: %v", err)
	// }
	// o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})
	//dat, _ := json.Marshal(data)
	readers := bytes.NewReader(data)
	// Upload an object with storage.Writer.
	wc := o.NewWriter(ctx)

	if _, err := io.Copy(wc, readers); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil
}
