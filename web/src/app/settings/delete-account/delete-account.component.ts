import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';


import { BackendService, AuthService } from "../../backend.service"
import { RValidators } from '../../forms/rvalidators';
import { ToastyService } from '../../toasty/toasty.module';

@Component({
  selector: 'rana-delete-account',
  templateUrl: './delete-account.component.html',
})
export class DeleteAccountComponent implements OnInit {
  form: FormGroup;

  constructor(
    private fb: FormBuilder,
    private backend: BackendService,
    private auth: AuthService,
    private toasty: ToastyService,
    private router: Router,
  ) { }

  ngOnInit() {
    this.form = this.fb.group({
      password: [null, RValidators.password],
    });
  }

  deleteAccount() {
    this.backend
      .deleteUser(this.auth.user, this.form.value.password)
      .subscribe(_ => {
        this.toasty.success('Account delete success');
        this.auth.unset();
        this.router.navigate(['/']);
      });
  }
}
