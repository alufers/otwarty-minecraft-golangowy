package omg

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type shaderProgram struct {
	vertexSource, fragmentSource string
	vertexShader, fragmentShader uint32
}

func newShaderProgram(vertexSource, fragmentSource string) *shaderProgram {
	sp := &shaderProgram{
		vertexSource:   vertexSource,
		fragmentSource: fragmentSource,
	}
	return sp
}

func (sp *shaderProgram) compile() error {
	sp.vertexShader = gl.CreateShader(gl.VERTEX_SHADER)
	vertexSourceStr, free := gl.Strs()
	gl.ShaderSource(sp.vertexShader, 1, vertexSourceStr, nil)
	free()
	gl.CompileShader(sp.vertexShader)
	var status int32
	gl.GetShaderiv(sp.vertexShader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(sp.vertexShader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(sp.vertexShader, logLength, nil, gl.Str(log))

		return fmt.Errorf("failed to compile vertex shader: %v", log)
	}

	sp.fragmentShader = gl.CreateShader(gl.FRAGMENT_SHADER)
	fragmentSourceStr, free := gl.Strs()
	gl.ShaderSource(sp.fragmentShader, 1, fragmentSourceStr, nil)
	free()
	gl.CompileShader(sp.fragmentShader)
	gl.GetShaderiv(sp.fragmentShader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(sp.fragmentShader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(sp.fragmentShader, logLength, nil, gl.Str(log))

		return fmt.Errorf("failed to compile fragment shader: %v", log)
	}
	return nil
}
