
-- sql2api.name: get_task_by_id
-- sql2api.description: usecase for task reading in context of TODO app
-- sql2api.http: GET
-- sql2api.async: yes
-- sql2api.parameters: id

SELECT * FROM todos WHERE id = 1;
	