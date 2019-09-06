package model

type (
	Guild struct {
		ID             string         `json:"id" db:"id"`
		Name           string         `json:"name" db:"name"`
		Icon           string         `json:"icon" db:"icon"`
		BotExists      bool           `db:"is_bot_exists"`
	}

	UserGuild struct {
		UserID      string `db:"user_id"`
		GuildID     string `json:"id" db:"guild_id"`
		IsOwner     bool   `json:"owner" db:"is_owner"`
		Permissions uint   `json:"permissions" db:"permissions"`
		CanInvite   bool   `db:"can_invite"`
	}
)

func (m *model) AddGuilds(guilds *[]Guild) (err error) {
	tx, err := m.db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	for _, guild := range *guilds {
		_, err = tx.Exec(`INSERT INTO 
						discord_guilds (id, name, icon, is_bot_exists) 
						VALUES (?, ?, ?, ?)
						ON DUPLICATE KEY UPDATE
						name = VALUES(name),
						icon = VALUES(icon),
						is_bot_exists = VALUES(is_bot_exists)`,
			guild.ID,
			guild.Name,
			guild.Icon,
			guild.BotExists)
		if err != nil {
			return
		}
	}

	return
}

func (m *model) AddUserGuild(userGuilds *[]UserGuild) (err error) {
	tx, err := m.db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	for _, ug := range *userGuilds {
		_, err = tx.Exec(`INSERT INTO users__discord_guilds
						VALUES (?, ?, ?, ?, ?)
						ON DUPLICATE KEY UPDATE
						is_owner = VALUES(is_owner),
						permissions = VALUES(permissions),
						can_invite = VALUES(can_invite)`,
			ug.UserID,
			ug.GuildID,
			ug.IsOwner,
			ug.Permissions,
			ug.CanInvite)
		if err != nil {
			return
		}
	}

	return
}

func (m *model) GetBotExistsGuilds() (*[]Guild, error) {
	guilds := &[]Guild{}
	if err := m.db.Select(guilds, `SELECT * FROM discord_guilds WHERE is_bot_exists = true`); err != nil {
		return nil, err
	}
	return guilds, nil
}

func (m *model) UpdateGuild(guild *Guild) (err error) {
	_, err = m.db.NamedQuery(`UPDATE discord_guilds SET 
							name=:name, 
							icon=:icon, 
							is_bot_exists=:is_bot_exists
							WHERE id=:id`, &guild)
	return err
}

func (m *model) UpdateGuildBotExists(id string, exists bool) (err error) {
	_, err = m.db.Exec(`UPDATE discord_guilds SET 
						is_bot_exists=? 
						WHERE id=?`, exists, id)
	return err
}
