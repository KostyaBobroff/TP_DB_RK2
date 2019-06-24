package database

import (
	"log"

	"github.com/Vlad104/TP_DB_RK2/models"
	"github.com/jackc/pgx"
)

const (
	createForumSQL = `
		INSERT INTO forums (slug, title, "user")
		VALUES ($1, $2, (
			SELECT nickname FROM users WHERE nickname = $3
		)) 
		RETURNING "user"
	`

	getForumSQL = `
		SELECT slug, title, "user", posts, threads
		FROM forums
		WHERE slug = $1
	`

	createForumThreadSQL = `
		INSERT INTO threads (author, created, message, title, slug, forum)
		VALUES ($1, $2, $3, $4, $5, (SELECT slug FROM forums WHERE slug = $6)) 
		RETURNING author, created, forum, id, message, title
	`

	getForumThreadsSinceSQL = `
		SELECT author, created, forum, id, message, slug, title, votes
		FROM threads
		WHERE forum = $1 AND created >= $2::TEXT::TIMESTAMPTZ
		ORDER BY created
		LIMIT $3::TEXT::INTEGER
	`
	getForumThreadsDescSinceSQL = `
		SELECT author, created, forum, id, message, slug, title, votes
		FROM threads
		WHERE forum = $1 AND created <= $2::TEXT::TIMESTAMPTZ
		ORDER BY created DESC
		LIMIT $3::TEXT::INTEGER
	`
	getForumThreadsSQL = `
		SELECT author, created, forum, id, message, slug, title, votes
		FROM threads
		WHERE forum = $1
		ORDER BY created
		LIMIT $2::TEXT::INTEGER
	`
	getForumThreadsDescSQL = `
		SELECT author, created, forum, id, message, slug, title, votes
		FROM threads
		WHERE forum = $1
		ORDER BY created DESC
		LIMIT $2::TEXT::INTEGER
	`
	getForumUsersSienceSQl = `
		SELECT forum_user, fullname, about, email
		FROM forum_users
		WHERE forum = $1
		AND LOWER(forum_user) > LOWER($2::TEXT)
		ORDER BY forum_user
		LIMIT $3::TEXT::INTEGER
	`
	getForumUsersDescSienceSQl = `
		SELECT forum_user, fullname, about, email
		FROM forum_users
		WHERE forum = $1
		AND LOWER(forum_user) < LOWER($2::TEXT)
		ORDER BY forum_user DESC
		LIMIT $3::TEXT::INTEGER
	`
	getForumUsersSQl = `
		SELECT forum_user, fullname, about, email
		FROM forum_users
		WHERE forum = $1
		ORDER BY forum_user
		LIMIT $2::TEXT::INTEGER
	`
	getForumUsersDescSQl = `
		SELECT forum_user, fullname, about, email
		FROM forum_users
		WHERE forum = $1
		ORDER BY forum_user DESC
		LIMIT $2::TEXT::INTEGER
	`
)

// /forum/create Создание форума
func CreateForumDB(f *models.Forum) (*models.Forum, error) {
	err := DB.pool.QueryRow(
		createForumSQL,
		&f.Slug,
		&f.Title,
		&f.User,
	).Scan(&f.User)

	switch ErrorCode(err) {
	case pgxOK:
		return f, nil
	case pgxErrUnique:
		forum, _ := GetForumDB(f.Slug)
		return forum, ForumIsExist
	case pgxErrNotNull:
		return nil, UserNotFound
	default:
		return nil, err
	}
}

// /forum/{slug}/details Получение информации о форуме
func GetForumDB(slug string) (*models.Forum, error) {
	f := models.Forum{}

	err := DB.pool.QueryRow(
		getForumSQL,
		slug,
	).Scan(
		&f.Slug,
		&f.Title,
		&f.User,
		&f.Posts,
		&f.Threads,
	)

	if err != nil {
		return nil, ForumNotFound
	}

	return &f, nil
}

// /forum/{slug}/create Создание ветки
func CreateForumThreadDB(t *models.Thread) (*models.Thread, error) {
	if t.Slug != "" {
		thread, err := GetThreadDB(t.Slug)
		if err == nil {
			return thread, ThreadIsExist
		}
	}

	err := DB.pool.QueryRow(
		createForumThreadSQL,
		&t.Author,
		&t.Created,
		&t.Message,
		&t.Title,
		&t.Slug,
		&t.Forum,
	).Scan(
		&t.Author,
		&t.Created,
		&t.Forum,
		&t.ID,
		&t.Message,
		&t.Title,
	)

	switch ErrorCode(err) {
	case pgxOK:
		return t, nil
	case pgxErrNotNull:
		return nil, ForumOrAuthorNotFound //UserNotFound
	case pgxErrForeignKey:
		return nil, ForumOrAuthorNotFound //ForumIsExist
	default:
		return nil, err
	}
}

var queryForumWithSience = map[string]string{
	"true":  getForumThreadsDescSinceSQL,
	"false": getForumThreadsSinceSQL,
}

var queryForumNoSience = map[string]string{
	"true":  getForumThreadsDescSQL,
	"false": getForumThreadsSQL,
}

// /forum/{slug}/threads Список ветвей обсужления форума
func GetForumThreadsDB(slug, limit, since, desc string) (*models.Threads, error) {
	var rows *pgx.Rows
	var err error

	if since != "" {
		query := queryForumWithSience[desc]
		rows, err = DB.pool.Query(query, slug, since, limit)
	} else {
		query := queryForumNoSience[desc]
		rows, err = DB.pool.Query(query, slug, limit)
	}
	defer rows.Close()

	if err != nil {
		return nil, ForumNotFound
	}

	threads := models.Threads{}
	for rows.Next() {
		t := models.Thread{}
		err = rows.Scan(
			&t.Author,
			&t.Created,
			&t.Forum,
			&t.ID,
			&t.Message,
			&t.Slug,
			&t.Title,
			&t.Votes,
		)
		threads = append(threads, &t)
	}

	if len(threads) == 0 {
		_, err := GetForumDB(slug)
		if err != nil {
			return nil, ForumNotFound
		}
	}
	return &threads, nil
}

var queryForumUserWithSience = map[string]string{
	"true":  getForumUsersDescSienceSQl,
	"false": getForumUsersSienceSQl,
}

var queryForumUserNoSience = map[string]string{
	"true":  getForumUsersDescSQl,
	"false": getForumUsersSQl,
}

// /forum/{slug}/users Пользователи данного форума
func GetForumUsersDB(slug, limit, since, desc string) (*models.Users, error) {
	var rows *pgx.Rows
	var err error

	if since != "" {
		query := queryForumUserWithSience[desc]
		rows, err = DB.pool.Query(query, slug, since, limit)
	} else {
		query := queryForumUserNoSience[desc]
		rows, err = DB.pool.Query(query, slug, limit)
	}
	defer rows.Close()

	if err != nil {
		log.Println(err)
		return nil, ForumNotFound
	}

	users := models.Users{}
	for rows.Next() {
		u := models.User{}
		err = rows.Scan(
			&u.Nickname,
			&u.Fullname,
			&u.About,
			&u.Email,
		)
		users = append(users, &u)
	}

	if len(users) == 0 {
		_, err := GetForumDB(slug)
		if err != nil {
			return nil, ForumNotFound
		}
	}
	return &users, nil
}
