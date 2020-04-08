import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PresosComponent } from './presos.component';

describe('PresosComponent', () => {
  let component: PresosComponent;
  let fixture: ComponentFixture<PresosComponent>;

  beforeEach(async(() => {
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
