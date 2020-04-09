import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { NavComponent } from './view/nav/nav.component';
import { ReposComponent } from './view/home/repos/repos.component';
import { HttpClientModule }    from '@angular/common/http';
import { PresosComponent } from './view/home/presos/presos.component';
import { PostsComponent } from './view/home/posts/posts.component';
import { HomeComponent } from './view/home/home.component';
import { AboutComponent } from './view/about/about.component';
import { BookComponent } from './view/book/book.component';
import { ResumeComponent } from './view/resume/resume.component';
import { ContactComponent } from './view/contact/contact.component';


@NgModule({
  declarations: [
    AppComponent,
    NavComponent,
    ReposComponent,
    PresosComponent,
    PostsComponent,
    HomeComponent,
    AboutComponent,
    BookComponent,
    ResumeComponent,
    ContactComponent
    
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
