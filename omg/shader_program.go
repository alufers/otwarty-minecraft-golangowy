package omg

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type shaderProgram struct {
	vertexSource, fragmentSource          string
	vertexShader, fragmentShader, program uint32
}

func newShaderProgram(vertexSource, fragmentSource string) *shaderProgram {
	sp := &shaderProgram{
		vertexSource:   vertexSource,
		fragmentSource: fragmentSource,
	}
	return sp
}

func (sp *shaderProgram) compile() error {
	var err error
	sp.vertexShader, err = sp.compileSingleShader(sp.vertexSource, gl.VERTEX_SHADER)

	if err != nil {
		return err
	}
	sp.fragmentShader, err = sp.compileSingleShader(sp.fragmentSource, gl.FRAGMENT_SHADER)

	if err != nil {
		return err
	}
	sp.program = gl.CreateProgram()
	gl.AttachShader(sp.program, sp.vertexShader)
	gl.AttachShader(sp.program, sp.fragmentShader)
	gl.LinkProgram(sp.program)
	gl.DeleteShader(sp.fragmentShader)
	gl.DeleteShader(sp.vertexShader)
	return nil
}

func (sp *shaderProgram) compileSingleShader(source string, shaderType uint32) (uint32, error) {
	sh := gl.CreateShader(shaderType)
	fragmentSourceStr, free := gl.Strs(source + "\x00")
	gl.ShaderSource(sh, 1, fragmentSourceStr, nil)
	free()
	gl.CompileShader(sh)
	var status int32
	gl.GetShaderiv(sh, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(sh, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(sh, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile fragment shader: %v", log)
	}

	return sh, nil
}

func (sp *shaderProgram) use() {
	gl.UseProgram(sp.program)
}
