package db

func (d *dbImpl) GetID(email string) (string, error) {
	var id string
	query := `SELECT id FROM signup WHERE email = $1`
	err := d.db.QueryRow(query, email).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}
func (d *dbImpl) GetUserByEmail(email string) (string, string, error) {
	var id, password string
	query := `SELECT id, password FROM signup WHERE email = $1`
	err := d.db.QueryRow(query, email).Scan(&id, &password)
	if err != nil {
		return "", "", err
	}
	return id, password, nil
}
