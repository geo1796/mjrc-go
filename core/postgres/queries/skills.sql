-- name: GetSkills :many
SELECT *
FROM app.skills;

-- name: CreateSkill :exec
INSERT INTO app.skills (name,
                        youtube_video_id,
                        is_video_landscape,
                        level,
                        categories,
                        prerequisites)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateSkill :one
UPDATE app.skills
SET name               = $2,
    youtube_video_id   = $3,
    is_video_landscape = $4,
    level              = $5,
    categories         = $6,
    prerequisites      = $7
WHERE id = $1 RETURNING id;

-- name: DeleteSkill :one
DELETE
FROM app.skills
WHERE id = $1 RETURNING id;

-- name: SkillsFingerprint :one
SELECT COUNT(*)::bigint AS cnt, (COALESCE(MAX(updated_at), 'epoch'::timestamptz)) ::timestamptz AS max_updated_at
FROM app.skills;