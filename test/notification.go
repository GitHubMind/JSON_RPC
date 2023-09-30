package test

//
//type NotificationTestService struct {
//	unsubscribed            chan string
//	gotHangSubscriptionReq  chan struct{}
//	unblockHangSubscription chan struct{}
//}
//
//func (s *NotificationTestService) Echo(i int) int {
//	return i
//}
//
//func (s *NotificationTestService) Unsubscribe(subid string) {
//	if s.unsubscribed != nil {
//		s.unsubscribed <- subid
//	}
//}
//
//func (s *NotificationTestService) SomeSubscription(ctx context.Context, n, val int) (*rpc.Subscription, error) {
//	notifier, supported := rpc.NotifierFromContext(ctx)
//	if !supported {
//		return nil, rpc.ErrNotificationsUnsupported
//	}
//
//	// By explicitly creating an subscription we make sure that the subscription id is send
//	// back to the client before the first subscription.Notify is called. Otherwise the
//	// events might be send before the response for the *_subscribe method.
//	subscription := notifier.CreateSubscription()
//	go func() {
//		for i := 0; i < n; i++ {
//			if err := notifier.Notify(subscription.ID, val+i); err != nil {
//				return
//			}
//		}
//		select {
//		case <-notifier.Closed():
//		case <-subscription.Err():
//		}
//		if s.unsubscribed != nil {
//			s.unsubscribed <- string(subscription.ID)
//		}
//	}()
//	return subscription, nil
//}
//
//// HangSubscription blocks on s.unblockHangSubscription before sending anything.
//func (s *NotificationTestService) HangSubscription(ctx context.Context, val int) (*rpc.Subscription, error) {
//	notifier, supported := rpc.NotifierFromContext(ctx)
//	if !supported {
//		return nil, rpc.ErrNotificationsUnsupported
//	}
//	s.gotHangSubscriptionReq <- struct{}{}
//	<-s.unblockHangSubscription
//	subscription := notifier.CreateSubscription()
//
//	go func() {
//		notifier.Notify(subscription.ID, val)
//	}()
//	return subscription, nil
//}
