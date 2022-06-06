package repositories

import (
	"api/src/models"
	"database/sql"
)

type Publications struct {
	db *sql.DB
}

// NewPublicationRepository create one repository of publication
func NewPublicationRepository(db *sql.DB) *Publications {
	return &Publications{db}
}

func (repository Publications) Create(Publication models.Publication) (uint64, error) {
	statement, err := repository.db.Prepare(
		"INSERT INTO publications (title, content, author_id) values (?, ?, ?)",
	)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(Publication.Title, Publication.Content, Publication.AuthorId)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(lastId), nil
}

func (repository Publications) FindById(publicationId uint64) (models.Publication, error) {
	line, err := repository.db.Query(`
		SELECT p.*, u.nick FROM publications
		p INNER JOIN users u 
		ON u.id = p.author_id WHERE p.id = ?`,
		publicationId,
	)
	if err != nil {
		return models.Publication{}, err
	}
	defer line.Close()

	var publication models.Publication

	if line.Next() {
		if err = line.Scan(
			&publication.ID,
			&publication.Title,
			&publication.Content,
			&publication.AuthorId,
			&publication.Likes,
			&publication.CreatedAt,
			&publication.AuthorNick,
		); err != nil {
			return models.Publication{}, err
		}
	}

	return publication, nil
}

func (repository Publications) Find(userId uint64) ([]models.Publication, error) {
	lines, err := repository.db.Query(`
		SELECT DISTINCT p.*, u.nick FROM publications p 
		INNER JOIN users u ON u.id = p.author_id
		INNER JOIN followers f on p.author_id = f.user_id
		WHERE u.id = ? or f.follower_id = ?
		ORDER BY 1 DESC`,
		userId, userId,
	)
	if err != nil {
		return nil, err
	}
	defer lines.Close()

	var publications []models.Publication

	for lines.Next() {

		var publication models.Publication

		if err = lines.Scan(
			&publication.ID,
			&publication.Title,
			&publication.Content,
			&publication.AuthorId,
			&publication.Likes,
			&publication.CreatedAt,
			&publication.AuthorNick,
		); err != nil {
			return nil, err
		}

		publications = append(publications, publication)

	}

	return publications, nil
}

func (repository Publications) Update(publicationId uint64, publication models.Publication) error {
	statement, err := repository.db.Prepare("UPDATE publications SET title = ?, content = ? WHERE id = ?")
	if err != nil {
		return err
	}

	defer statement.Close()

	if _, err = statement.Exec(publication.Title, publication.Content, publicationId); err != nil {
		return err
	}

	return nil
}

func (repository Publications) Delete(publicationId uint64) error {
	statement, err := repository.db.Prepare("DELETE FROM publications WHERE id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(publicationId); err != nil {
		return err
	}

	return nil
}
