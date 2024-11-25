
-- sql2api.name: edit_task_by_id
-- sql2api.description: usecase for task updating in context of TODO app
-- sql2api.http: PUT
-- sql2api.async: yes
-- sql2api.parameters: auto

UPDATE todos
SET task = 'Comprar leite e pão', completed = TRUE
WHERE id = 1;
	