package mysql

import (
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strings"
)

type CreateReplicaUserSettings struct {
	Force    bool
	DryRun   bool
	Schema   string
	Username string
	Password string
}

type step struct {
	statement          string
	runningDescription string
	errorDescription   string
}

func (md *MysqlDB) CreateReplicaUser(settings CreateReplicaUserSettings) error {
	var err error

	sb := sqlbuilder.Select("user").From("mysql.user")
	sb.Where(sb.Equal("user", settings.Username))
	sb2 := sqlbuilder.Buildf("SELECT EXISTS(%v)", sb)
	sql_, args_ := sb2.Build()

	if settings.Force {
		sql_ = fmt.Sprintf("DROP USER IF EXISTS %s", settings.Username)
		if settings.DryRun {
			log.Info().Str("sql", sql_).Msg("Force deletion of user")
		} else {
			_, err = md.Db.Exec(sql_)
			if err != nil {
				return errors.Wrap(err, "Could not delete user")
			}
			log.Info().Str("username", settings.Username).Msg("Deleting user")
		}
	} else {
		if settings.DryRun {
			log.Info().Str("sql", sql_).Interface("args", args_).Msg("Checking if user exists")
		} else {
			var exists bool
			err = md.Db.QueryRow(sql_, args_...).Scan(&exists)
			if err != nil {
				return errors.Wrap(err, "Could not check if user exists")
			}
			if exists {
				log.Fatal().Str("username", settings.Username).Msg("User already exists")
			}
		}
	}

	statements := []step{
		{"CREATE USER ${User}", "Creating user", "Could not create user"},
		{"SET PASSWORD FOR ${User} = PASSWORD('${Password}')", "Setting password", "Could not set password"},
		{"GRANT ALL ON ${Schema}.* TO ${User}", "Granting privileges", "Could not grant privileges"},
		{"GRANT RELOAD ON *.* TO ${User}", "Granting reload privileges", "Could not grant reload privileges"},
		{"GRANT REPLICATION CLIENT ON *.* TO ${User}", "Granting replication privileges", "Could not grant replication privileges"},
		{"GRANT REPLICATION SLAVE ON *.* TO ${User}", "Granting replication slave privileges", "Could not grant replication slave privileges"},
		{"FLUSH PRIVILEGES", "Flushing privileges", "Could not flush privileges"},
	}
	replacer := strings.NewReplacer("${User}", settings.Username, "${Password}", settings.Password, "${Schema}", settings.Schema)

	for _, s := range statements {
		sql_ := replacer.Replace(s.statement)

		if settings.DryRun {
			log.Info().Str("sql", sql_).Msg(s.runningDescription)
		} else {
			_, err = md.Db.Exec(sql_)
			if err != nil {
				return errors.Wrap(err, s.errorDescription)
			}
			log.Info().Str("username", settings.Username).Msg(s.runningDescription)
		}
	}

	return nil
}
