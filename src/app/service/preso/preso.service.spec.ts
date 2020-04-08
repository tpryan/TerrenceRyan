import { TestBed } from '@angular/core/testing';

import { PresoService } from './preso.service';

describe('PresoService', () => {
  let service: PresoService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(PresoService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
