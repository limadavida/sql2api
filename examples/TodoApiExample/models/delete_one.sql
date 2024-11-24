-- name: del_task_by_id
-- description: usecase for task deleting in context of TODO app
-- async: yes
-- parameters: id

DELETE FROM todos WHERE id = 1;
