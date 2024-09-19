select ut.user_id   as user_id,
       ut.ticket_id as ticket_id
from user_tickets ut
where ut.user_id = $1;