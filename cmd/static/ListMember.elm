module ListMember exposing (Model, Msg, init, view, update)

import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onCheck)

-- MODEL

type alias Model =
  { id: Int
    , name: String
    , desc: String
    , selected: Bool
  }

init: Int -> String -> String -> Model
init pid pname pdesc = { id = pid, name = pname, desc = pdesc, selected = False }

-- UPDATE

type Msg = Selected Bool

update: Msg -> Model -> Model
update msg model =
  case msg of
    Selected b ->
      { model | selected = b }

-- SUBSCRIPTIONS

-- VIEW

view: Model -> Html Msg
view model =
  div []
    [
      input [ type' "checkbox", onCheck Selected ] []
      , text model.name
      , text model.desc
    ]
