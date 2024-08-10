package dto


type CategoriaDto struct {
	Nombre string `json:"nombre"`
}

type GenericoDto struct {
	Estado  string `json:"estado"`
	Mensaje string `json:"mensaje"`
}

type ProductoDto struct {
	Nombre      string  `json:"nombre"`
	Precio      float64 `json:"precio"`
	Stock       int     `json:"stock"`
	Descripcion string  `json:"descripcion"`
	CategoriaID string  `json:"categoria_id"`
}

type UsuarioDto struct {
	Nombre   string `json:"nombre"`
	Correo   string `json:"correo"`
	Password string `json:"password"`
	Telefono string `json:"telefono"`
}

type LoginDto struct {
	Correo   string `json:"correo"`
	Password string `json:"password"`
}

type LoginRespuestaDto struct {
	Nombre string `json:"nombre"`
	Token  string `json:"token"`
}
