package utils

import (
	"strconv"
	"strings"
)

// Color - Renkler
type Colors struct {
	Success int
	Warning int
	Error   int
	Info    int
}

// Color - Global Renkler (int tipindeler)
var ColorsInt *Colors

// Color - Renk modülünü kurar.
func InitColors() error {
	// YAML Tarafındaki kodlar HEX durumunda , tip çevirme işlemi için geçici bir yapı
	raw := &struct {
		Success string `yaml:"success"`
		Warning string `yaml:"warning"`
		Error   string `yaml:"error"`
		Info    string `yaml:"info"`
	}{}

	// Renkleri yamldan oku
	if err := ReadYaml("config/core/color.yaml", raw); err != nil {
		return err
	}

	// Renkleri Global Renklere ata
	ColorsInt = &Colors{
		Success: hexToInt(raw.Success),
		Error:   hexToInt(raw.Error),
		Warning: hexToInt(raw.Warning),
		Info:    hexToInt(raw.Info),
	}

	return nil
}

// Color - Hex kodunda olan rengi Int e çevirir
func hexToInt(hex string) int {
	clean := strings.TrimPrefix(hex, "#")
	val, err := strconv.ParseInt(clean, 16, 32)
	if err != nil {
		return 0
	}
	return int(val)
}
