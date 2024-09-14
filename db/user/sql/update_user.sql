update users
set email   = :email,
    name    = :name,
    surname = :surname
where id = :id;