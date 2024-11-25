
-- sql2api.name: del_task_by_id
-- sql2api.description: usecase for task deleting in context of TODO app
-- sql2api.http: DELETE
-- sql2api.async: yes
-- sql2api.parameters: id

DELETE FROM todos WHERE id = 1;
		