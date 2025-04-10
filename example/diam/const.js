export const cmd = {
    Accounting:          271,
    CreditControl:       272,
}

export const cmdFlag = {
    Request:             0x80, // request bit
    Proxiable:           0x40, // proxiable bit
    Error:               0x20, // error bit
    Retransmit:          0x10, // re-transmitted bit
}

export const app = {
    Accounting:          3,
    ChargingControl:     4,
}

export const flag = {
    V:                   0x80, // vendor bit
    M:                   0x40, // mandatory bit
    P:                   0x20, // private bit
}
export const code = {
    CalledStationId:               30,
    DestinationHost:               293,
    DestinationRealm:              283,
    OriginHost:                    264,
    OriginRealm:                   296,
    OriginStateId:                 278,
    SessionId:                     263,
    CCRequestNumber:               415,
    CCRequestType:                 416,
    ResultCode:                    268,
    SubscriptionId:                443,
    SubscriptionIdData:            444,
    SubscriptionIdType:            450,
    ServiceInformation:            873,
    PSInformation:                 874,
    // Extra
    EventNameCode:                 13001,
}

export const vendor = {
    TGPP:                          10415,
    Matrixxsoftware:               35838,
}

