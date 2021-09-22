import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { PresosComponent } from './presos.component';

describe('PresosComponent', () => {
  let component: PresosComponent;
  let fixture: ComponentFixture<PresosComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ PresosComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PresosComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
