import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { AboutComponent } from './view/about/about.component';
import { HomeComponent } from './view/home/home.component';

import { BookComponent } from './view/book/book.component';
import { ResumeComponent } from './view/resume/resume.component';
import { ContactComponent } from './view/contact/contact.component';



const routes: Routes = [
  { path: '', redirectTo: '/home', pathMatch: 'full' },
  { path: 'home', component: HomeComponent },
  { path: 'about', component: AboutComponent },
  { path: 'book', component: BookComponent },
  { path: 'resume', component: ResumeComponent },
  { path: 'contact', component: ContactComponent },


  
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
