import TwitterList
import Html exposing (..)
import Html.App as Html
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Http
import Task
import Json.Decode as Json


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
  | FetchFail String

update: Msg -> Model -> (Model, Cmd Msg)
update msg model =
    case msg of
      GetLists ->
        (model, getTwitterLists)
      UpdateList m ->
        (model, Cmd.none)
      FetchSuccess plists ->
        ({ model | lists = plists }, Cmd.none)
      FetchFail _ ->
        (model, Cmd.none)


getTwitterLists: Cmd Msg
getTwitterLists =
  let url = "http://localhosts/lists"
  in Task.perform FetchFail FetchSuccess (Http.get (Json.list Json.string) url)

decodeLists: Json.Decoder String
decodeLists =
  Json.at ["data", "image_url"] Json.string

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
