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
