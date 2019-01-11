module Modals.Remove exposing (Model, Msg, newModel, show, subscriptions, update, view)

import Bootstrap.Alert as Alert
import Bootstrap.Button as Button
import Bootstrap.Form.Input as Input
import Bootstrap.Grid as Grid
import Bootstrap.Grid.Col as Col
import Bootstrap.Grid.Row as Row
import Bootstrap.Modal as Modal
import Bootstrap.Progress as Progress
import Browser
import Browser.Navigation as Nav
import File
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Http
import Json.Decode as D
import Json.Encode as E
import List
import Ls
import Url
import Util



-- TODO: Handle case where the dir already exist.
--       Warn and ask to overwrite?


type State
    = Ready
    | Fail String


type alias Model =
    { state : State
    , modal : Modal.Visibility
    , alert : Alert.Visibility
    }


type Msg
    = RemoveAll (List String)
    | ModalShow
    | GotResponse (Result Http.Error String)
    | AnimateModal Modal.Visibility
    | AlertMsg Alert.Visibility
    | ModalClose



-- INIT


newModel : Model
newModel =
    { state = Ready
    , modal = Modal.hidden
    , alert = Alert.shown
    }



-- UPDATE


type alias Query =
    { paths : List String
    }


encode : Query -> E.Value
encode q =
    E.object
        [ ( "paths", E.list E.string q.paths ) ]


decode : D.Decoder String
decode =
    D.field "message" D.string


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        RemoveAll paths ->
            ( model
            , Http.post
                { url = "/api/v0/remove"
                , body = Http.jsonBody <| encode <| Query paths
                , expect = Http.expectJson GotResponse decode
                }
            )

        GotResponse result ->
            case result of
                Ok _ ->
                    -- New list model means also new checked entries.
                    ( { model | state = Ready, modal = Modal.hidden }, Cmd.none )

                Err err ->
                    ( { model | state = Fail <| Util.httpErrorToString err }, Cmd.none )

        AnimateModal visibility ->
            ( { model | modal = visibility }, Cmd.none )

        ModalShow ->
            ( { model | modal = Modal.shown }, Cmd.none )

        ModalClose ->
            ( { model | modal = Modal.hidden, state = Ready }, Cmd.none )

        AlertMsg vis ->
            ( { model | alert = vis }, Cmd.none )



-- VIEW


viewRemoveContent : Model -> Ls.Model -> List (Grid.Column Msg)
viewRemoveContent model lsModel =
    [ Grid.col [ Col.xs12 ]
        [ case model.state of
            Ready ->
                text ("Remove the " ++ String.fromInt (Ls.nSelectedItems lsModel) ++ " selected items")

            Fail message ->
                Util.buildAlert model.alert AlertMsg True "Oh no!" ("Could not remove directory: " ++ message)
        ]
    ]


view : Model -> Ls.Model -> Html Msg
view model lsModel =
    Modal.config ModalClose
        |> Modal.large
        |> Modal.withAnimation AnimateModal
        |> Modal.h5 [] [ text "Really remove?" ]
        |> Modal.body []
            [ Grid.containerFluid []
                [ Grid.row [] (viewRemoveContent model lsModel) ]
            ]
        |> Modal.footer []
            [ Button.button
                [ Button.danger
                , Button.attrs
                    [ onClick <| RemoveAll <| Ls.selectedPaths lsModel
                    , disabled
                        (case model.state of
                            Fail _ ->
                                True

                            _ ->
                                False
                        )
                    ]
                ]
                [ text "Remove" ]
            , Button.button
                [ Button.outlinePrimary
                , Button.attrs [ onClick <| AnimateModal Modal.hiddenAnimated ]
                ]
                [ text "Cancel" ]
            ]
        |> Modal.view model.modal


show : Msg
show =
    ModalShow



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.batch
        [ Modal.subscriptions model.modal AnimateModal
        , Alert.subscriptions model.alert AlertMsg
        ]
