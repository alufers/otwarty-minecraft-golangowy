package omg

import "github.com/go-gl/gl/v4.1-core/gl"

type mesh struct {
	shaderProgram *shaderProgram
	vertices      []float32
	vbo           uint32
	vao           uint32
}

func (m *mesh) bindBuffers() {
	gl.GenBuffers(1, &m.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(m.vertices), gl.Ptr(m.vertices), gl.STATIC_DRAW)
	gl.GenVertexArrays(1, &m.vao)
	gl.BindVertexArray(m.vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

}

func (m *mesh) draw() {
	m.shaderProgram.use()
	gl.BindVertexArray(m.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(m.vertices)/3))
}
