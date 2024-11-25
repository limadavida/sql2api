
-- sql2api.name: create_task
-- sql2api.description: usecase for task creation in context of TODO app
-- sql2api.http: POST
-- sql2api.async: yes
-- sql2api.parameters: magic

INSERT INTO todos (task, completed) VALUES ('Comprar leite', 1);
	