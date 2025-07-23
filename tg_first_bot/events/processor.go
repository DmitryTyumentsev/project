package events

//
//func (p *EventProcessor) process(event Event) error {
//	switch event.Type {
//	case Message:
//		doCMD(Message)
//		//if err := processMessage(event); err != nil {
//		//	return errors.WrapIfErr("failed processMessage", err)
//		//}
//		//if err := p.client.SendMessage(event.ChatID, ""); err != nil {
//		//	return errors.WrapIfErr("failed to send msg", err)
//		//}
//	default:
//		doCMD(UnknownType)
//		//if err := p.client.SendMessage(event.ChatID, ""); err != nil {
//		//	return errors.WrapIfErr("failed to send msg", err)
//		//}
//	}
//	return nil
//}
//
//func processMessage(event Event) error {
//	if err := doCMD(event); err != nil {
//		return errors.WrapIfErr("failed processMessage", err)
//	}
//	return nil
//}
