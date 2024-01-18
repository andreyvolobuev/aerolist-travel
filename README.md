# Aerolist Travel module #

I want to re-implement all the logic of Aerolist from Python to Golang and then migrate it.  
I don't know how much time will it take and will I actually finish it or not but I'm willing to at least give it a try.
My first module of choice to re-implement is travel which is in my opinion is the main module of the app.

## API ##
trips/
    GET:
        ?arr=<id:optional>  # get trips from given city
        ?dep=<id:optional>  # get trips to given city
        ?date=<date:optional>  # get trips within given bounds
    POST:  # create new trip

trips/<id>/
    GET:  # detailed info about trip
    PUT:  # update data in trip
    DELETE:  # delete a trip

cities/
    GET:
        ?q=<name:optional>  # for autocomplete

cities/<id>/
    GET:  # detailed info about city

countries/
    GET:
        ?q=<name:optional>  # for autocomplete

countries/<id>/
    GET:  # detailed info about a country

trip_view_requests/
    GET:  # get list of trip view requests
    POST:  # create new trip view request

trip_view_requests/<id>/
    GET:  # detailed info about trip view request
    PUT:  # update data in trip view request
    DELETE:  # delete a trip view request
