import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MakeCaseComponent } from './make-case.component';

describe('MakeCaseComponent', () => {
  let component: MakeCaseComponent;
  let fixture: ComponentFixture<MakeCaseComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MakeCaseComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MakeCaseComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
