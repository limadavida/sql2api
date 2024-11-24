-- name: edit_task_by_id
-- description: usecase for task updating in context of TODO app
-- async: yes
-- parameters: auto

UPDATE todos
SET task = 'Comprar leite e pão', completed = TRUE
WHERE id = 2;
