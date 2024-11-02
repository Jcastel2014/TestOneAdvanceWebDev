SELECT R.id, P.name, R.rating, R.helpful_count, R.comment, R.created_at, R.updated_at 
FROM reviews AS R
INNER JOIN products AS P ON P.id = R.product_id
WHERE P.id = 1 AND R.product_id = 6