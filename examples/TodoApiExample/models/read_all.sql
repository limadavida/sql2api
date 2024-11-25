
-- sql2api.name: get_all_tasks
-- sql2api.description: usecase for task reading in context of TODO app
-- sql2api.http: GET
-- sql2api.async: yes
-- sql2api.parameters: max=10

SELECT * FROM todos;
	