import TwitterList
import Html exposing (..)
import Html.App as Html
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Http
import Task
import Json.Decode as Json exposing (..)
import Debug


main =
  Html.program
    { init = init
    , view = view
    , update = update
    , subscriptions = subscriptions
    }

-- MODEL

type alias Model =
  {
    lists: List TwitterList.Model
  }

init: (Model, Cmd Msg)
init = (Model [], getTwitterLists)
  

-- UPDATE

type Msg = UpdateList TwitterList.Msg
  | GetLists
  | FetchSuccess (List TwitterList.Model)
  | FetchFail Http.Error

update: Msg -> Model -> (Model, Cmd Msg)
update msg model =
    case msg of
      GetLists ->
        (model, getTwitterLists)
      UpdateList m ->
        (model, Cmd.none)
      FetchSuccess plists ->
        Debug.log "Entering FetchSuccess"
        ({ model | lists = plists }, Cmd.none)
      FetchFail error ->
        Debug.log (toString error)
        (model, Cmd.none)


getTwitterLists: Cmd Msg
getTwitterLists =
  --let url = "http://localhost:8080/lists"
  let url = "/lists"
  in Task.perform FetchFail FetchSuccess (Http.get decodeLists url)

decodeLists: Json.Decoder (List TwitterList.Model)
decodeLists =
  Debug.log "Entering decodeLists"
  Json.at ["Lists"] (Json.list decodeList)

decodeList: Json.Decoder TwitterList.Model
decodeList =
  Json.object3 TwitterList.Model ("id" := Json.int) ("name" := Json.string) (Json.null [])

-- VIEW

view: Model -> Html Msg
view model =
  div [] (List.map viewTwitterList model.lists)

viewTwitterList: TwitterList.Model -> Html Msg
viewTwitterList model = 
  Html.map UpdateList (TwitterList.view model)

-- SUBSCRIPTIONS

subscriptions : Model -> Sub Msg
subscriptions model =
  Sub.none
