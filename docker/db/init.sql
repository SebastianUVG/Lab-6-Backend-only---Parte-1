-- Crear la tabla matches
CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,
    home_team VARCHAR(255) NOT NULL,
    away_team VARCHAR(255) NOT NULL,
    match_date DATE NOT NULL
);

-- Insertar datos iniciales (opcional)
INSERT INTO matches (home_team, away_team, match_date) 
VALUES ('Barcelona', 'Real Madrid', '2025-04-01')
ON CONFLICT DO NOTHING;