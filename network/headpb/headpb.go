package headpb

// type Processor struct {
// 	msgInfo map[uint32]*MsgInfo
// 	msgID   map[reflect.Type]uint32
// }

// type MsgInfo struct {
// 	msgType       reflect.Type
// 	msgRouter     *chanrpc.Server
// 	msgHandler    MsgHandler
// 	msgRawHandler MsgHandler
// }

// type MsgHandler func([]interface{})

// type MsgRaw struct {
// 	msgID      uint32
// 	msgRawData []byte
// 	head       *HeadUnit
// 	body       interface{}
// }

// func NewProcessor() *Processor {
// 	p := new(Processor)
// 	p.msgInfo = make(map[uint32]*MsgInfo)
// 	p.msgID = make(map[reflect.Type]uint32)
// 	return p
// }

// // It's dangerous to call the method on routing or marshaling (unmarshaling)
// func (p *Processor) Register(msgID uint32, msg interface{}) *Processor {
// 	msgType := reflect.TypeOf(msg)
// 	if msgType == nil || msgType.Kind() != reflect.Ptr {
// 		log.Fatalf("message pointer required")
// 	}
// 	if _, ok := p.msgInfo[msgID]; ok {
// 		log.Fatalf("message %v is already registered", msgID)
// 	}
// 	i := new(MsgInfo)
// 	i.msgType = msgType
// 	p.msgInfo[msgID] = i
// 	p.msgID[msgType] = msgID
// 	return p
// }

// // It's dangerous to call the method on routing or marshaling (unmarshaling)
// func (p *Processor) SetRouter(msgID uint32, msgRouter *chanrpc.Server) *Processor {
// 	i, ok := p.msgInfo[msgID]
// 	if !ok {
// 		log.Fatalf("msgId: %d not registered", msgID)
// 	}
// 	i.msgRouter = msgRouter
// 	return p
// }

// // It's dangerous to call the method on routing or marshaling (unmarshaling)
// func (p *Processor) SetRawHandler(msgID uint32, msgRawHandler MsgHandler) *Processor {
// 	i, ok := p.msgInfo[msgID]
// 	if !ok {
// 		log.Fatalf("message %d not registered", msgID)
// 	}
// 	i.msgRawHandler = msgRawHandler
// 	return p
// }

// // goroutine safe
// func (p *Processor) Unmarshal(data []byte) (interface{}, error) {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			log.Errorf("%v", r)
// 		}
// 	}()
// 	var err error
// 	nUnit := &HeadUnit{}
// 	nUnit.Decode(data[0:24])
// 	i, ok := p.msgInfo[nUnit.MsgID]
// 	if !ok {
// 		err = fmt.Errorf("message %v not registered", nUnit.MsgID)
// 		log.Errorf("%v", err)
// 		return nil, err
// 	}
// 	// 构建RequestMessage结构
// 	msgOut := reflect.New(i.msgType.Elem()).Interface()
// 	log.Debugf("[recv] msgId: %d, cextral: %d", nUnit.MsgID, nUnit.ClientExtral)

// 	err = proto.UnmarshalMerge(data[24:], msgOut.(proto.Message))
// 	if err == nil {
// 		return &MsgRaw{msgID: nUnit.MsgID, msgRawData: data, head: nUnit, body: msgOut}, nil
// 	} else {
// 		log.Errorf("unmarshalMerge error:%v", err)
// 		return nil, err
// 	}
// }

// // goroutine safe
// func (p *Processor) Route(msg interface{}, userData interface{}) error {
// 	if msgRaw, ok := msg.(*MsgRaw); ok {
// 		i, ok := p.msgInfo[msgRaw.msgID]
// 		if !ok {
// 			return fmt.Errorf("message %v not registered", msgRaw.msgID)
// 		}
// 		if i.msgRawHandler != nil {
// 			i.msgRawHandler([]interface{}{msgRaw.msgID, msgRaw.msgRawData, msgRaw.head, msgRaw.body, userData})
// 		} else {
// 			msgType := reflect.TypeOf(msgRaw.body)
// 			if i.msgRouter != nil {
// 				i.msgRouter.Go(msgType, msgRaw.body, userData, msgRaw.head)
// 			}
// 		}
// 		return nil
// 	}
// 	return errors.New("route msg invalidate!!!")
// }

// // goroutine safe
// func (p *Processor) Marshal(msg interface{}) ([][]byte, error) {
// 	var nUnit = &HeadUnit{From: 12}
// 	var body interface{}
// 	if msgRaw, ok := msg.(*MsgRaw); ok {
// 		nUnit.ClientExtral = msgRaw.head.ClientExtral
// 		body = msgRaw.body
// 	} else {
// 		body = msg
// 	}
// 	msgType := reflect.TypeOf(body)
// 	if msgType == nil || msgType.Kind() != reflect.Ptr {
// 		return nil, errors.New("marshal fail, message type dont exsit")
// 	}
// 	var msgID uint32
// 	var ok bool
// 	if msgID, ok = p.msgID[msgType]; !ok {
// 		return nil, errors.New("marshal fail, not registered message")
// 	}
// 	nUnit.MsgID = msgID
// 	data, err := proto.Marshal(body.(proto.Message))
// 	if err != nil {
// 		log.Errorf("proto Marshal Error:%v", err)
// 	}
// 	nUnit.MsgLen = uint32(len(data))
// 	nHeadBytes, _ := nUnit.Encode()

// 	if nUnit.MsgID != 1 {
// 		log.Debugf("[send] msgID:%d,len:%d,data:%+v,extral:%v", msgID, nUnit.MsgLen, body, nUnit.ClientExtral)
// 	}
// 	return [][]byte{nHeadBytes, data}, nil
// }
