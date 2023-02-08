package image_merge

import (
	"image"
	"image/png"
	"log"
	"os"
	"path"
	"path/filepath"
)

type MergeEngine struct {
	AssetPath       string
	OutputPath      string
	ResultPath      string
	BackgroundImage string
	Layers          []*Grid
	Merger          *MergeImage
	WriteFunc       WriterFunc
}

type EngineOpt func(engine *MergeEngine)

func WithDefaultWriter() EngineOpt {
	return func(engine *MergeEngine) {
		engine.WriteFunc = func(rgba *image.RGBA, path string) error {
			engine.ResultPath = filepath.Join(engine.OutputPath, path)
			file, err := os.Create(engine.ResultPath)
			log.Println("create file:", err)
			if err != nil {
				return err
			}
			return png.Encode(file, rgba)
		}
	}
}

func WithWriter(writer WriterFunc) EngineOpt {
	return func(engine *MergeEngine) {
		engine.WriteFunc = writer
	}
}

func WithOutputPath(path string) EngineOpt {
	return func(engine *MergeEngine) {
		engine.OutputPath = path
	}
}

func WithReader(reader ReadFunc) EngineOpt {
	return func(engine *MergeEngine) {
		engine.Merger.ReadFunc = reader
	}
}

func NewMergeEngine(assetPath string, outputPath string, opts ...EngineOpt) *MergeEngine {
	grid := make([]*Grid, 1)
	m := &MergeEngine{
		AssetPath:  assetPath,
		OutputPath: outputPath,
		Merger:     New(grid, 1, 1),
	}
	defaultOpts := []EngineOpt{WithDefaultWriter()}
	if len(opts) == 0 {
		for _, opt := range defaultOpts {
			opt(m)
		}
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *MergeEngine) Background(imagePath string) *MergeEngine {
	b := &Grid{
		ImageFilePath: path.Join(m.AssetPath, imagePath),
		Grids:         make([]*Grid, 1),
	}
	m.Merger.Grids[0] = b
	return m
}

func (m *MergeEngine) AddLayer(imagePath string) *MergeEngine {
	m.Layers = append(m.Layers, &Grid{
		ImageFilePath: path.Join(m.AssetPath, imagePath),
	})
	return m
}
func (m *MergeEngine) FromLayers(imagePath ...string) *MergeEngine {
	for _, p := range imagePath {
		m.Layers = append(m.Layers, &Grid{
			ImageFilePath: path.Join(m.AssetPath, p),
		})
	}
	return m
}

func (m *MergeEngine) Merge(output string) error {
	m.Merger.Grids[0].Grids = m.Layers
	rgba, err := m.Merger.Merge()
	if err != nil {
		return err
	}
	return m.WriteFunc(rgba, filepath.Join(m.OutputPath, output))
}
