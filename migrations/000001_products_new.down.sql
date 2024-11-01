DROP TABLE users;


	Update images
	SET image_link =$6
	WHERE id = (SELECT image_id FROM products WHERE id = $5);