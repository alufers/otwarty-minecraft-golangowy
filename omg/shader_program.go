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

// compile compiles the shader and returns an error if it fails
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
	gl.DeleteShader(sp.fragmentShader) // when the shaders are linked they can be safely deleted from opengl state
	gl.DeleteShader(sp.vertexShader)
	return nil
}

// compileSingleShader compiles a single shader (vertex or fragment) with error handling
func (sp *shaderProgram) compileSingleShader(source string, shaderType uint32) (uint32, error) {
	sh := gl.CreateShader(shaderType)
	fragmentSourceStr, free := gl.Strs(source + "\x00") // strings in c are zero terminated
	gl.ShaderSource(sh, 1, fragmentSourceStr, nil)
	free() // we need to free c strings
	gl.CompileShader(sh)
	var status int32
	gl.GetShaderiv(sh, gl.COMPILE_STATUS, &status) // get the status of compilation
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(sh, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1)) // allocate memory for the log with one byte left for the zero terminator
		gl.GetShaderInfoLog(sh, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile fragment shader: %v", log)
	}

	return sh, nil
}

// use enables this shader in opengl for next draw calls
func (sp *shaderProgram) use() {
	gl.UseProgram(sp.program)
}
