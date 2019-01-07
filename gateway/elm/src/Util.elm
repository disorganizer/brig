module Util exposing (basename, formatLastModified, monthToInt, splitPath, urlToPath)

import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Time
import Url


boolToStr : Bool -> String
boolToStr b =
    case b of
        True ->
            "yes"

        False ->
            "no"


monthToInt : Time.Month -> Int
monthToInt month =
    -- This feels stupid.
    case month of
        Time.Jan ->
            1

        Time.Feb ->
            2

        Time.Mar ->
            3

        Time.Apr ->
            4

        Time.May ->
            5

        Time.Jun ->
            6

        Time.Jul ->
            7

        Time.Aug ->
            8

        Time.Sep ->
            9

        Time.Oct ->
            10

        Time.Nov ->
            11

        Time.Dec ->
            12


formatLastModified : Time.Zone -> Time.Posix -> String -> Html.Html msg
formatLastModified z t owner =
    let
        timestamp =
            String.join " "
                -- Day portion:
                [ String.join
                    "/"
                    [ Time.toDay z t |> String.fromInt
                    , Time.toMonth z t |> monthToInt |> String.fromInt
                    , Time.toYear z t |> String.fromInt
                    ]

                -- Time portion:
                , String.join ":"
                    [ Time.toHour z t |> String.fromInt |> String.padLeft 2 '0'
                    , Time.toMinute z t |> String.fromInt |> String.padLeft 2 '0'
                    , Time.toSecond z t |> String.fromInt |> String.padLeft 2 '0'
                    ]
                ]
    in
    p [] [ text timestamp, span [ class "text-muted" ] [ text " by " ], text owner ]


splitPath : String -> List String
splitPath path =
    List.filter (\s -> String.length s > 0) (String.split "/" path)


urlToPath : Url.Url -> String
urlToPath url =
    let
        decodeUrlPart =
            \e ->
                case Url.percentDecode e of
                    Just val ->
                        val

                    Nothing ->
                        ""
    in
    case splitPath url.path of
        [] ->
            "/"

        _ :: xs ->
            "/" ++ String.join "/" (List.map decodeUrlPart xs)


basename : String -> String
basename path =
    let
        splitUrl =
            List.reverse (splitPath path)
    in
    case splitUrl of
        [] ->
            "/"

        x :: _ ->
            x