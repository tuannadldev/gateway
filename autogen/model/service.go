package model

type Service struct {
    Name           string
    Path           string
    Package        string
    Apis           []*Api
}
