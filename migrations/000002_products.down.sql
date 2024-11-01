SELECT 
    P.id, 
    P.name, 
    P.description, 
    P.category, 
    I.image_url, 
    P.average_rating, 
    P.created_at, 
    P.updated_at
FROM 
    products AS P
INNER JOIN 
    images AS I ON P.image_id = I.id
WHERE 
    P.id = 3;
