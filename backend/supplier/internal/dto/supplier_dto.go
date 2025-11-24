package dto

type CreateSupplierDTO struct {
	KodePemasok  string `json:"kode_pemasok"`
	NamaPemasok  string `json:"nama_pemasok" validate:"required,min=3"`
	KontakPerson string `json:"kontak_person"`
	Email        string `json:"email" validate:"omitempty,email"`
	Telepon      string `json:"telepon" validate:"required"`
	Alamat       string `json:"alamat"`
}

type UpdateSupplierDTO struct {
	KodePemasok  string `json:"kode_pemasok"`
	NamaPemasok  string `json:"nama_pemasok" validate:"omitempty,min=3"`
	KontakPerson string `json:"kontak_person"`
	Email        string `json:"email" validate:"omitempty,email"`
	Telepon      string `json:"telepon"`
	Alamat       string `json:"alamat"`
}
