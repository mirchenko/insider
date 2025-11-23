do $$
  begin
    for i in 0..1000
      loop
        insert into messages(status, content, phone_number)
        values ((select s from unnest(array['pending', 'sent', 'delivered', 'failed']) as s order by random() limit 1)::messages_type,
                'Content ' || (select substr(md5(random()::text), 1, 10)),
                (select pn from unnest(array['+34111111111', '+34222222222', '+34333333333']) as pn order by random() limit 1));
      end loop;
  end;
  $$ language plpgsql;
