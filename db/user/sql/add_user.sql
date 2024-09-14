insert into users (email, name, surname)
values (:email, :name, :surname)
returning id;