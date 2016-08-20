module TwitterList exposing(Model, init, Msg, update, view)

import ListMember
import Html exposing (..)
import Html.Attributes exposing (..)

-- MODEL
type alias Model =
  {
    id: Int
    , name: String
    , members: List ListMember.Model
  }

init: Int -> String -> Model
init id name = Model id name []

-- UPDATE

type Msg = Add ListMember.Model | Remove ListMember.Model

update: Msg -> Model -> Model
update msg model =
  case msg of
    Add member ->
      { model | members = model.members ++ [member] }
    Remove member ->
      model
      -- TODO
 
-- VIEW

view: Model -> Html Msg
view model =
  div []
  [
    a [ href "" ] [ text model.name ] 
  ]
