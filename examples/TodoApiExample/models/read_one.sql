-- name: get_task_by_id
-- description: usecase for task reading in context of TODO app
-- async: yes
-- parameters: id

SELECT * FROM todos WHERE id = 1;