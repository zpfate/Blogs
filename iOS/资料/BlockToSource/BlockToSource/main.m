//
//  main.m
//  BlockToSource
//
//  Created by ws on 2021/5/28.
//


#import <Foundation/Foundation.h>

int main() {
    // self
     NSObject *self = [NSObject new];
    __weak NSObject *weakSelf = self;
    void(^block)(void) = ^{
        weakSelf;
    };
}

