-- name: GetAllChatRoom
SELECT
  cr.id,
  cr.name,
  cr.private,
  cr.description,
  cr.userid,
  COUNT(gm.id) AS members
FROM
  chatroom cr
  LEFT JOIN groupMembers gm ON cr.name = gm.roomName
WHERE
  cr.private = false
  OR cr.userid = $1
GROUP BY
  cr.id;

-- name: GetAllJoinRoom
SELECT
  cr.id,
  cr.name,
  cr.description,
  cr.private,
  cr.userId,
  COUNT(gm.id) AS members
FROM
  chatRoom cr
  LEFT JOIN groupMembers gm ON cr.name = gm.roomName
WHERE
  gm.userId = $1
GROUP BY
  cr.id;

-- name: GetOldMessages
SELECT
  *
FROM
  (
    SELECT
      m.id,
      m.userId,
      u.username,
      m.content,
      m.roomName,
      m.createdAt AS messageCreatedAt,
      m.updatedAt AS messageUpdatedAt,
      f.id AS fileId,
      f.publicId,
      f.secureUrl,
      f.format,
      f.resourceType,
      f.size,
      f.width,
      f.height,
      f.originalFilename,
      f.createdAt AS fileCreatedAt,
      f.updatedAt AS fileUpdatedAt
    FROM
      messages m
      JOIN users u ON m.userId = u.id
      LEFT JOIN files f ON m.fileId = f.id
    WHERE
      m.roomName = $1
    ORDER BY
      m.createdAt DESC
    LIMIT
      $2
  ) subquery
ORDER BY
  messageCreatedAt ASC;

-- name: GetPrivateRoomUsingCode
SELECT
  cr.id,
  cr.name,
  cr.private,
  cr.description,
  cr.userid,
  COUNT(gm.id) AS members
FROM
  chatroom cr
  LEFT JOIN groupMembers gm ON cr.name = gm.roomName
WHERE
  cr.private = true
  AND cr.code = $1
GROUP BY
  cr.id;

--name: GetMessageWithFile
SELECT
  m.id,
  m.userId,
  m.roomName,
  m.content,
  m.createdAt,
  m.updatedAt,
  f.id AS file_id,
  f.name AS file_name,
  f.size AS file_size
FROM
  messages m
  LEFT JOIN files f ON m.fileId = f.id
WHERE
  m.roomName = $1
ORDER BY
  m.createdAt DESC
LIMIT
  $2;
