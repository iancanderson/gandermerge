//
//  ViewController.m
//  spookypaths
//
//  Created by Ian Anderson on 11/4/22.
//

#import "ViewController.h"
@import GoogleMobileAds;

@interface ViewController ()<GADFullScreenContentDelegate>

@property(nonatomic, strong) GADBannerView *bannerView;
@property(nonatomic, strong) GADRewardedAd *rewardedAd;

@end

@implementation ViewController

- (void)loadRewardedAd {
  GADRequest *request = [GADRequest request];
  [GADRewardedAd
       loadWithAdUnitID:@"ca-app-pub-3940256099942544/1712485313"
                request:request
      completionHandler:^(GADRewardedAd *ad, NSError *error) {
        if (error) {
          NSLog(@"Rewarded ad failed to load with error: %@", [error localizedDescription]);
          return;
        }
        self.rewardedAd = ad;
        NSLog(@"Rewarded ad loaded.");
        self.rewardedAd.fullScreenContentDelegate = self;
      
        [self.rewardedAd presentFromRootViewController:self
                                    userDidEarnRewardHandler:^{
                                    GADAdReward *reward =
                                        self.rewardedAd.adReward;
                                    // TODO: Reward the user!
                                  }];
      
  }];
}

/// Tells the delegate that the ad failed to present full screen content.
- (void)ad:(nonnull id<GADFullScreenPresentingAd>)ad
    didFailToPresentFullScreenContentWithError:(nonnull NSError *)error {
    NSLog(@"Ad did fail to present full screen content.");
}

/// Tells the delegate that the ad will present full screen content.
- (void)adWillPresentFullScreenContent:(nonnull id<GADFullScreenPresentingAd>)ad {
    NSLog(@"Ad will present full screen content.");
}

/// Tells the delegate that the ad dismissed full screen content.
- (void)adDidDismissFullScreenContent:(nonnull id<GADFullScreenPresentingAd>)ad {
    NSLog(@"Ad did dismiss full screen content.");
}

- (void)viewDidLoad {
    [super viewDidLoad];
    // Do any additional setup after loading the view.
    
    // In this case, we instantiate the banner with desired ad size.
    self.bannerView = [[GADBannerView alloc]
        initWithAdSize:GADAdSizeBanner];
    
    [self addBannerViewToView:self.bannerView];

    self.bannerView.adUnitID = @"ca-app-pub-3940256099942544/2934735716";
    self.bannerView.rootViewController = self;
    
    [self.bannerView loadRequest:[GADRequest request]];

    [self loadRewardedAd];
}

- (void)addBannerViewToView:(UIView *)bannerView {
  bannerView.translatesAutoresizingMaskIntoConstraints = NO;
  [self.view addSubview:bannerView];
  [self.view addConstraints:@[
    [NSLayoutConstraint constraintWithItem:bannerView
                               attribute:NSLayoutAttributeBottom
                               relatedBy:NSLayoutRelationEqual
                                  toItem:self.bottomLayoutGuide
                               attribute:NSLayoutAttributeTop
                              multiplier:1
                                constant:0],
    [NSLayoutConstraint constraintWithItem:bannerView
                               attribute:NSLayoutAttributeCenterX
                               relatedBy:NSLayoutRelationEqual
                                  toItem:self.view
                               attribute:NSLayoutAttributeCenterX
                              multiplier:1
                                constant:0]
                                ]];
}

@end
