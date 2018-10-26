SELECT trips.route_id, trips.service_id, trips.trip_id, trips.trip_headsign, trips.trip_short_name 
    FROM trips trips, calendar_dates cal
    WHERE trips.trip_short_name = '75' 
        AND trips.service_id = cal.service_id
        AND cal.date :: TEXT=TO_CHAR(NOW() :: DATE, 'yyyymmdd')
    LIMIT 1;


SELECT stops.stop_id, stops.stop_name, stops.stop_lat, stops.stop_lon, 
        EXTRACT(EPOCH FROM st.arrival_time - to_char(now() at time zone 'UTC-02', 'HH24:MI:SS') :: TIME) as arrival_sec, 
        st.stop_sequence 
    FROM stops stops, stop_times st
    WHERE st.trip_id='20749000010031'
        AND st.stop_id = stops.stop_id
    ;
