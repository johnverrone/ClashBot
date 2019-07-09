package bot_test

import (
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/johnverrone/clashbot/pkg/bot"
	"github.com/johnverrone/clashbot/pkg/chat/chatfakes"
	"github.com/johnverrone/clashbot/pkg/clash"
	"github.com/johnverrone/clashbot/pkg/clash/clashfakes"
)

var (
	fakeClashClient *clashfakes.FakeClient
	fakeChatClient  *chatfakes.FakeClient
	err             error
	prevState       *PrevState
)

var _ = Describe("Bot", func() {
	Describe("RunBotLogic", func() {
		BeforeEach(func() {
			fakeClashClient = &clashfakes.FakeClient{}
			fakeChatClient = &chatfakes.FakeClient{}
		})

		Context("there is an error getting the war", func() {
			BeforeEach(func() {
				fakeClashClient.GetWarReturns(clash.CurrentWar{}, errors.New("some error"))
				err = RunBotLogic(fakeClashClient, fakeChatClient, &PrevState{})
			})

			It("logs an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("getting the war is successful", func() {
			Context("and the war just started", func() {
				BeforeEach(func() {
					fakeClashClient.GetWarReturns(clash.CurrentWar{
						State: "inWar",
					}, nil)
					err = RunBotLogic(fakeClashClient, fakeChatClient, &PrevState{War: "preparation"})
				})

				It("calls SendMessage with war results", func() {
					Expect(err).ToNot(HaveOccurred())
					Expect(fakeChatClient.SendMessageCallCount()).To(Equal(1))
				})
			})

			Context("and the war just ended", func() {
				BeforeEach(func() {
					fakeClashClient.GetWarReturns(clash.CurrentWar{
						State: "warEnded",
					}, nil)
					err = RunBotLogic(fakeClashClient, fakeChatClient, &PrevState{War: "inWar"})
				})

				It("calls SendMessage with war results", func() {
					Expect(err).ToNot(HaveOccurred())
					Expect(fakeChatClient.SendMessageCallCount()).To(Equal(1))
				})
			})

			Context("and we are currently in a war", func() {
				BeforeEach(func() {
					fakeClashClient.CheckForAttackUpdatesReturns("")
					layout := "20060102T150405.000Z"
					endTime := time.Now().UTC().Add(2 * time.Minute).Format(layout)
					prevState = &PrevState{SentWarReminder: false}
					members := []clash.ClanWarMember{
						{
							Name: "john",
							Attacks: []clash.ClanWarAttack{
								{
									AttackerTag: "john",
								},
								{
									AttackerTag: "john",
								},
							},
						},
						{
							Name: "alex",
							Attacks: []clash.ClanWarAttack{
								{
									AttackerTag: "alex",
								},
							},
						},
						{
							Name:    "newb",
							Attacks: []clash.ClanWarAttack{},
						},
					}
					fakeClashClient.GetWarReturns(clash.CurrentWar{
						State:   "inWar",
						EndTime: endTime,
						Clan: clash.WarClan{
							Members: members,
						},
					}, nil)
					err = RunBotLogic(fakeClashClient, fakeChatClient, prevState)
				})

				It("calls SendMessage", func() {
					Expect(err).ToNot(HaveOccurred())
					Expect(fakeChatClient.SendMessageCallCount()).To(Equal(1))
					Expect(prevState.SentWarReminder).To(Equal(true))
				})
			})
		})
	})
})
